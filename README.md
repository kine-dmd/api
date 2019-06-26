# Kine-DMD API [![Build Status](https://travis-ci.org/kine-dmd/api.svg?branch=master)](https://travis-ci.org/kine-dmd/api)
Go API for deployment to AWS ECS and for local running on server. Used to match data sent directly a watch with a universally unique identifier to a patient ID and transmit the data onwards for processing.
![image](https://user-images.githubusercontent.com/26333869/60197045-7e87dd80-9836-11e9-8bfe-6f94d39d7b65.png)

## AWS ECS
The API is primarily designed to be run on AWS ECS. To this effect, it is designed to be stateless to allow for multiple containers to be run in parallel.
![ecsContainers](https://user-images.githubusercontent.com/26333869/60198515-b5132780-9839-11e9-9e5c-19ea27df7781.png)



## Local run
The API can also be run on a local server. When this is done, the API will initialise a different watch database and data writer on start up.
![image](https://user-images.githubusercontent.com/26333869/60198766-4a162080-983a-11e9-9220-289c6d86d49f.png)

To run the container, the below folder structure must first be established with folders created manually for each patient. Failure to create folders for user data will result in API errors. In addition, the users CSV file and appropriate public and private keys for the web server must be initialised as follows:

![image](https://user-images.githubusercontent.com/26333869/60197407-39b07680-9837-11e9-93a3-c23a3fba9b74.png)

Once this is done, tag the built container as kine-dmd-api and run the following command. This maps the folders outside the container to folders inside the container and uses the docker daemon to restart the container if it crashes or the underlying OS is rebooted. It also runs the process in the background so a user can safely exit their session whilst leaving the API running.
```
docker run -p 0.0.0.0:3000:443 -e "kine_dmd_api_location=local" -v $PWD/certs/:/certs -v $PWD/watch_db/:/watch_db -v $PWD/data:/data --restart always -d kine-dmd-api
```


