Build the project and use the binary main.exe with the BurpSuite's Logger++ CSV File

1. Go to Burp Suite's Logger++ 
2. Select the requests and click Export Entities as CSV
3. Now select option to export request as Base64 only
4. Now use the main.exe with -filepath flag to pass the csv file
5. A swagger.yaml will be generated 
6. Go to editor.swagger.io
7. Click import file
8. Import the swagger.yaml file
9. Now the swagger documentation of all the apis will be generated
