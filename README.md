This repository contains unofficial patterns, sample code, or tools to help developers build more effectively with [Fauna][fauna]. All [Fauna Labs][fauna-labs] repositories are provided “as-is” and without support. By using this repository or its contents, you agree that this repository may never be officially supported and moved to the [Fauna organization][fauna-organization].

---

# Node Express with Typescript Fly.io starter for Fauna

*[Fauna](https://fauna.com/) is a distributed relational database with a document data model. Delivered as an API, Fauna is automatically configured – out of the box – as a three replica database with active-active write capability, making it a powerful complement to [Fly.io](https://fly.io/) in serving low latency reads and writes for dynamic global applications.*

This starter kit provides a sample Fauna integration with Go/Gin framework, and packaged to run on Fly.io. Whether you plan to deploy on Fly or not, this sample should nevertheless provide a great Go/Gin + Fauna example.
 
---

## Prerequisites
* Go 1.19 or greater
* [flyctl](https://fly.io/docs/hands-on/install-flyctl/)

## Create a Fauna database and generate an access key

* [Signup](https://dashboard.fauna.com/register) for a Fauna account if you don't have one already.
* Create a database and database access token according to [these instructions](https://docs.fauna.com/fauna/current/get_started/client_quick_start?lang=go).
* Copy the `.env.example` file (in the root of this project) into a new file named `.env` and populate the variable FAUNA_SECRET_KEY with the database access token from above:
  ```
  export FAUNA_SECRET_KEY="xxxxxxxx-xxxxxxxxx"
  ```

## Test locally

Source `.env`

```
source .env
```

__Let's load some sample data so that we can do a little manual test__. There is a node script in the `/scripts` folder for this. Run:
```
npm install --prefix ./scripts && node ./scripts/sample/load.js
```
At this point you should check your Fauna dashboard shell and you'll notice a few sample collections and data loaded from the script above. Now you can tidy: 

```
go mod tidy
```

...then start the server

```
go run main.go
```

Navigate to `http://localhost:8080/read`

This should perform a read from the Fauna database you created and populated with sample data, per the setup steps above.


## Deploy to Fly.io

To launch the app on fly, run `fly launch --no-deploy` in the root directory of this project.
You will be prompted for a couple things:

* `? Would you like to copy its configuration to the new app? (y/N)` **Choose y (Yes)**
* `? Do you want to tweak these settings before proceeding? (y/N)` **Choose N (No)**

```bash
% fly launch --no-deploy
An existing fly.toml file was found for app go-gin-fauna-starter
? Would you like to copy its configuration to the new app? Yes
Using build strategies '[a buildpack]'. Remove [build] from fly.toml to force a rescan
Creating app in <folder>/go-gin-fly-io-starter
We're about to launch your app on Fly.io. Here's what you're getting:

Organization: <your org name>           (fly launch defaults to the personal org)
Name:         go-gin-fauna-starter      (from your fly.toml)
Region:       San Jose, California (US) (from your fly.toml)
App Machines: shared-cpu-1x, 1GB RAM    (most apps need about 1GB of RAM)
Postgres:     <none>                    (not requested)
Redis:        <none>                    (not requested)

? Do you want to tweak these settings before proceeding? No
Created app 'go-gin-fauna-starter' in organization 'personal'
Admin URL: https://fly.io/apps/go-gin-fauna-starter
Hostname: go-gin-fauna-starter.fly.dev
Wrote config file fly.toml
Validating <folder>/go-gin-fly-io-starter/fly.toml
Platform: machines
✓ Configuration is valid
Your app is ready! Deploy with `flyctl deploy`
```

Before deploying, you must set the Secrets value for FAUNA_SECRET_KEY: 
```
fly secrets set FAUNA_SECRET_KEY="<fauna secret key>"
```

Now you can deploy:
```
fly deploy
```

Once the application has been deployed, you can find out more about its deployment. 
```
fly status
```

Browse to your newly deployed application with the `fly open` command.
```
% fly open

Opening https://<a-new-app-name>.fly.dev
```

## Scale your deployment to match the Fauna footprint

When you create a database in Fauna, it is automatically configured – out of the box – as a three replica database with active-active write capability (For example, if you created a database in the "US Region Group", there will be 3 replicas of the database across the United States). Thus, a good way to take advantage of this architecture is to deploy on 3 Fly.io regions as well, as close as possible to the database replicas.

<img src="docs/img/RG.png" alt="Region Groups" width="370">

Currently, Fauna provides 2 choices of Regions Groups, US and EU. The table below lists the Fly regions that are closest to the Fauna replicas of each respective region group:

| Fauna Region Group | Deploy on Fly Regions |
|--------------------|-----------------------|
| EU                 | lhr, arn, fra         |
| US                 | sjc, ord, iad         |

For example, let's say you created your Fauna database in the US Region Group. This starter kit is provided with a default [`fly.toml`](./fly.toml) file with `primary_region` set to `sjc`. Unless you edited this value, deploying this starter kit leaves you with your Fly app in `sjc`. To take full advantage of Fauna’s distributed footprint, add additional Fly machines in the other 2 regions closest to the Fauna replicas by running this command: 

```
fly scale count 2 --region ord,iad
```

Then, run `fly scale show` to see where your app’s Machines are running. For example:

```
$ fly scale show

VM Resources for app: my-app-name

Groups
NAME    COUNT   KIND    CPUS    MEMORY  REGIONS
app     3       shared  1       256 MB  iad,ord,sjc
```

There is nothing else that needs to be updated in the code or Fauna configuration, because Fauna automatically routes requests to the closest replica based on latency and availability. 


[fauna]: https://www.fauna.com/
[fauna-labs]: https://github.com/fauna-labs
[fauna-organization]: https://github.com/fauna
