#!/bin/bash

echo "Applying aws-auth ConfigMap with correct github-ci user mapping..."

# Use default profile (which has proper permissions)
export AWS_PROFILE=default

# Update kubeconfig with default profile
aws eks update-kubeconfig --name simple-bank --region ap-south-1 --profile default

# Apply the aws-auth ConfigMap
kubectl apply -f eks/aws-auth.yaml

echo "aws-auth ConfigMap applied successfully!"
echo "Now the github-ci user should be able to authenticate to the cluster." 