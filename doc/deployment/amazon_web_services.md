# Amazon Web Services

Below, we show how to deploy Pachyderm on AWS in a couple of different ways:

1. [By manually deploying Kubernetes and Pachyderm.](amazon_web_services.html#manual-pachyderm-deploy)
2. [By executing a one shot deploy script that will both deploy Kubernetes and Pachyderm.](amazon_web_services.html#one-shot-script)

If you already have a Kubernetes deployment or would like to customize the types of instances, size of volumes, etc. in your Kubernetes cluster, you should follow option (1).  If you just want a quick deploy to experiment with Pachyderm in AWS or would just like to use our default configuration, you might want to try option (2)


## Production Deployment

Note - for production deployments we recommend setting up AWS CloudFront. AWS puts S3 rate limits in place that can limit the data throughput for your cluster, and CloudFront helps mitigate this issue.

[Follow the instructions here to deploy a Pachyderm cluster with CloudFront](./aws_cloudfront.html)

## Prerequisites

- [AWS CLI](https://aws.amazon.com/cli/) - have it installed and have your [AWS credentials](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) configured.
- [kubectl](https://kubernetes.io/docs/user-guide/prereqs/)
- [kops](https://github.com/kubernetes/kops/blob/master/docs/install.md)

## Manual Pachyderm Deploy

### Deploy Kubernetes

The easiest way to install Kubernetes on AWS is with kops. Kubenetes has provided a [step by step guide](https://github.com/kubernetes/kops/blob/master/docs/aws.md) for the deploy.  Please follow [this guide](https://github.com/kubernetes/kops/blob/master/docs/aws.md) to deploy Kubernetes on AWS.  

Once, you have a Kubernetes cluster up and running in AWS, you should be able to see the following output from `kubectl`:

```shell
$ kubectl get all
NAME             CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
svc/kubernetes   10.0.0.1     <none>        443/TCP   22s
```

### Deploy Pachyderm

To deploy Pachyderm we will need to:

1. Install the `pachctl` CLI tool,
2. Add some storage resources on AWS, 
3. Deploy Pachyderm on top of the storage resources.

#### Install `pachctl`

To deploy and interact with Pachyderm, you will need `pachctl`, a command-line utility used for Pachyderm. To install `pachctl` run one of the following:


```shell
# For OSX:
$ brew tap pachyderm/tap && brew install pachyderm/tap/pachctl@1.5

# For Linux (64 bit):
$ curl -o /tmp/pachctl.deb -L https://github.com/pachyderm/pachyderm/releases/download/v1.5.0/pachctl_1.5.0_amd64.deb && sudo dpkg -i /tmp/pachctl.deb
```

You can try running `pachctl version` to check that this worked correctly, but Pachyderm itself isn't deployed yet so you won't get a `pachd` version.

```sh
$ pachctl version
COMPONENT           VERSION
pachctl             1.4.6
pachd               (version unknown) : error connecting to pachd server at address (0.0.0.0:30650): context deadline exceeded.
```

#### Set up the Storage Resources

Pachyderm needs an [S3 bucket](https://aws.amazon.com/documentation/s3/), and a [persistent disk](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumes.html) (EBS) to function correctly.

Here are the environmental variables you should set up to create these resources:

```shell
$ kubectl cluster-info
  Kubernetes master is running at https://1.2.3.4
  ...
$ KUBECTLFLAGS="-s [The public IP of the Kubernetes master. e.g. 1.2.3.4]"

# BUCKET_NAME needs to be globally unique across the entire AWS region
$ BUCKET_NAME=[The name of the S3 bucket where your data will be stored]

# We recommend between 1 and 10 GB. This stores PFS metadata. For reference 1GB
# should work for 1000 commits on 1000 files.
$ STORAGE_SIZE=[the size of the EBS volume that you are going to create, in GBs. e.g. "10"]

$ AWS_REGION=[the AWS region of your Kubernetes cluster. e.g. "us-west-2" (not us-west-2a)]

$ AWS_AVAILABILITY_ZONE=[the AWS availability zone of your Kubernetes cluster. e.g. "us-west-2a"]

```

Then to actually create the resources, you can run:
```
$ aws s3api create-bucket --bucket ${BUCKET_NAME} --region ${AWS_REGION} --create-bucket-configuration LocationConstraint=${AWS_REGION}

$ aws ec2 create-volume --size ${STORAGE_SIZE} --region ${AWS_REGION} --availability-zone ${AWS_AVAILABILITY_ZONE} --volume-type gp2
```

Record the "volume-id" that is output (e.g. "vol-8050b807") from the above `create-volume` command as shown below (you can also view it in the aws console or with  `aws ec2 describe-volumes`):

```shell
$ STORAGE_NAME=<volume id>
```

Now, as a sanity check, you should be able to see the bucket and the EBS volume that are just created:

```shell
aws s3api list-buckets --query 'Buckets[].Name'
aws ec2 describe-volumes --query 'Volumes[].VolumeId'
```

#### Deploy Pachyderm

When you installed kops, you should have created a dedicated IAM user (see [here](https://github.com/kubernetes/kops/blob/master/docs/aws.md#aws) for details).  To deploy Pachyderm you will need to export these credentials to the following environmental variables:

```sh
$ AWS_ACCESS_KEY_ID=[access key ID]

$ AWS_SECRET_ACCESS_KEY=[secret access key]
```

Run the following command to deploy your Pachyderm cluster:

```shell
$ pachctl deploy amazon ${BUCKET_NAME} ${AWS_ACCESS_KEY_ID} ${AWS_SECRET_ACCESS_KEY} " " ${AWS_REGION} ${STORAGE_SIZE} --static-etcd-volume=${STORAGE_NAME}
```

(Note, the `" "` in the deploy command is for an optional temporary AWS token, if you are just experimenting with a deploy.  Such a token should NOT be used for a production deploy).  It may take a few minutes for the pachd nodes to be running because it's pulling containers from DockerHub. You can see the cluster status by using:

```sh
$ kubectl get all
NAME                        READY     STATUS    RESTARTS   AGE
po/etcd-4197107720-br61m    1/1       Running   0          8m
po/pachd-3548222380-s086m   1/1       Running   2          8m

NAME             CLUSTER-IP     EXTERNAL-IP   PORT(S)                       AGE
svc/etcd         10.111.11.36   <nodes>       2379:32379/TCP                8m
svc/kubernetes   10.96.0.1      <none>        443/TCP                       10m
svc/pachd        10.97.116.5    <nodes>       650:30650/TCP,651:30651/TCP   8m

NAME           DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deploy/etcd    1         1         1            1           8m
deploy/pachd   1         1         1            1           8m

NAME                  DESIRED   CURRENT   READY     AGE
rs/etcd-4197107720    1         1         1         8m
rs/pachd-3548222380   1         1         1         8m
```

Note: If you see a few restarts on the pachd nodes, that's totally ok. That simply means that Kubernetes tried to bring up those containers before etcd was ready so it restarted them.

Finally, we need to set up forward a port so that pachctl can talk to the cluster.

```sh
# Forward the ports. We background this process because it blocks.
$ pachctl port-forward &
```

And you're done! You can test to make sure the cluster is working by trying `pachctl version` or even creating a new repo.
```sh
$ pachctl version
COMPONENT           VERSION
pachctl             1.4.6
pachd               1.4.6
```

## One Shot Script

### Install additional prerequisites

This scripted deploy requires a couple of prerequisites in addition to the ones listed under [Prerequisites](amazon_web_services.md#prerequisites):

- [jq](https://stedolan.github.io/jq/download/)
- [uuid](http://man7.org/linux/man-pages/man1/uuidgen.1.html)

### Run the deploy script

Once you have the prerequisites mentioned above, download and run our AWS deploy script by running:

```
curl -o aws.sh https://raw.githubusercontent.com/pachyderm/pachyderm/master/etc/deploy/aws.sh
chmod +x aws.sh
sudo -E ./aws.sh
```

This script will use kops to deploy Kubernetes and Pachyderm in AWS.  The script will ask you for your AWS credentials, region preference, etc.  If you would like to customize the number of nodes in the cluster, node types, etc., you can open up the deploy script and modify the respective fields.

The script will take a few minutes, and Pachyderm will take an addition couple of minutes to spin up.  Once it is up, `kubectl get all` should return something like:

```
NAME                        READY     STATUS    RESTARTS   AGE
po/etcd-4197107720-br61m    1/1       Running   0          8m
po/pachd-3548222380-s086m   1/1       Running   2          8m

NAME             CLUSTER-IP     EXTERNAL-IP   PORT(S)                       AGE
svc/etcd         10.111.11.36   <nodes>       2379:32379/TCP                8m
svc/kubernetes   10.96.0.1      <none>        443/TCP                       10m
svc/pachd        10.97.116.5    <nodes>       650:30650/TCP,651:30651/TCP   8m

NAME           DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deploy/etcd    1         1         1            1           8m
deploy/pachd   1         1         1            1           8m

NAME                  DESIRED   CURRENT   READY     AGE
rs/etcd-4197107720    1         1         1         8m
rs/pachd-3548222380   1         1         1         8m
```

### Connect `pachctl`

Finally, we need to set up forward a port so that pachctl can talk to the cluster.

```sh
# Forward the ports. We background this process because it blocks.
$ pachctl port-forward &
```

And you're done! You can test to make sure the cluster is working by trying `pachctl version`:
```sh
$ pachctl version
COMPONENT           VERSION
pachctl             1.4.6
pachd               1.4.6
```
