Steps:

1. Set up the home-assistant instance

    - open Terminal
    - `docker-compose up`
    - go to `URL:8123`
    - Make account - not https so don't use an important password 
    (location and name doesn't matter)

    - click finish

2. get an api key

    - click on name (bottom left)
    - Scroll to bottom of the page
    - Create a token under Long-Lived Access Tokens
    - Give it a name and click ok
    - Save the token as we will use it later
    - Click ok

3. Save api key
    - Rename terraform.tfvars.example to terraform.tfvars
    - In the file replace YOUR-TOKEN in bearer_token= "YOUR-TOKEN" with the token you just made (If you have lost your token do step 2 again)
