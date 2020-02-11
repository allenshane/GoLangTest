# GoLangTest
Test application with golang

The app performs the follow actions:
- Reads a configuration file in JSON format, gets all required configuration. The configuration file for the application should be passed as a command line argument.
- Copies all directories mentioned in the configuration to the destination folder mentioned in the configuration.
- Dumps and backs up all the MySQL databases mentioned in the configuration to the same folder.
- Compresses the entire folder into a single compressed file
- The resulting compressed file is uploaded to an AWS S3 bucket, the details for which is taken from the configuration file.
