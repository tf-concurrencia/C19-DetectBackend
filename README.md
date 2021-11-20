# C19-DetectBackend
Proyecto backend

In root of each MS, execute commnand in terminal:
-for ReadDataset MS
docker build -t read-dataset-svc
docker run -d -p 8081:8081 load-dataset-svc
-for TrainModel MS
docker build -t read-train-model-svc
docker run -d -p 8082:8082 train-model-svc
-for PredictCovid MS
docker build -t predict-covid-svc
docker run -d -p 8083:8083 predict-covid-svc

try in Browser:
http://host.docker.internal:8081/load-dataset
or 
http://host.docker.internal:8082/train-model

User Postman or ThunderCliente to send an array of boolean in JSON format to 3rd MS
http://host.docker.internal:8083/predict-covid

{
    "inputs": [1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0]
}

returns:
{
  "rpta": "1"
}
