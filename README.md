# Employee API - a POC on Google Cloud Run and Terraform

## Purpose

This POC is to explore how Google Cloud Run can be used as a solution for some simple microservices. 

Terraform is used to deploy all the required resources to support the idea of infrastructure as code

## Local testing

Only Docker engine is required. docker compose is used here.

## Manual deployment

1. Enable APIs
	```bash
	gcloud services enable sqladmin.googleapis.com
	gcloud services enable sql-component.googleapis.com
	```
2. Create a Cloud SQL instance named `sql-db`
	```bash
	gcloud sql instances create sql-db \
	--tier db-f1-micro \
	--database-version MYSQL_8_0 \
	--region us-central1
	```
3. Install cloud sql proxy
	```bash
	curl -o cloud-sql-proxy https://storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/v2.4.0/cloud-sql-proxy.linux.amd64
	chmod +x cloud-sql-proxy
	sudo mv cloud-sql-proxy /usr/local/bin 
	```
4. Install MySQL
	```bash
	sudo apt-get install mysql-client
	```
5. Get the instance connection name
	```bash
	gcloud sql instances describe sql-db|grep -i connection | awk '{print $2}'
	```
6. Start the cloud sql proxy to connect to `sql-db`
	```bash
	cloud-sql-proxy --port 3306 ibcwe-event-layer-f3ccf6d9:us-central1:sql-db
	Connect with the mysql client
	mysql -u root --host 127.0.0.1
	```
7. Test the connection
	```bash
	mysql> show databases;
	+--------------------+
	| Database           |
	+--------------------+
	| information_schema |
	| mysql              |
	| performance_schema |
	| sys                |
	+--------------------+
	4 rows in set (0.00 sec)
	```
8. Create the `hr` database
	```bash
	mysql> create database todo;
	```
9. Create the the todos table and add a few rows into it from the schema.sql
	```bash
	mysql -u root --host 127.0.0.1 todo < schema.sql
	```
10. Enable Artifact Registry in your project
	```bash
	gcloud services enable artifactregistry.googleapis.com
	```
11. Create a repository named `cloud-run-try` on `us-central1`. 
	```bash
	gcloud artifacts repositories create \
	--location us-central1 \
	--repository-format docker \
	cloud-run-try
	```
12. Build the docker image.
	```bash
	docker build -t \
	us-central1-docker.pkg.dev/<gcp-project-id>/cloud-run-try/employee .
	```
13. Set up credentials to access the repo
	```bash
	gcloud auth configure-docker us-central1-docker.pkg.dev
	```
14. Push the image
	```bash
	docker push \
	us-central1-docker.pkg.dev/<gcp-project-id>/cloud-run-try/employee
	```
15. Deploy the Cloud Run service
	```bash
	gcloud run deploy employee-api \
	--image us-central1-docker.pkg.dev/<gcp-project-id>/cloud-run-try/employee \
	--allow-unauthenticated`
	```
16. Check the service is deployed successfully
	```bash
	gcloud run service list
	```

## Deploying with Terraform

1. Initialization
	```bash
	terraform init
	```
2. Creating a plan
	```bash
	terraform plan -out tfplan
	```
3. Applying the plan
	```bash
	terraform apply
	```

## Testing

Get the service endpoint using the following command:

```bash
$ gcloud run services list
   SERVICE       REGION       URL   
✔  employee-api  us-central1  https://employee-api-oy6beuif2a-uc.a.run.app  
```

Hit the endpoint

```bash
curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
https://employee-api-oy6beuif2a-uc.a.run.app/api/employees
```

## Cleanup

Delete all the resources
```bash
terraform destroy
```

Delete the container image
```bash
gcloud artifacts packages delete employee --repository=cloud-run-try --location=us-central1
```
