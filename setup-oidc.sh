#!/bin/bash

# Setup OIDC Provider for GitHub Actions
echo "Setting up OIDC Provider for GitHub Actions..."

# Create OIDC Provider
aws iam create-open-id-connect-provider \
    --url https://token.actions.githubusercontent.com \
    --client-id-list sts.amazonaws.com \
    --thumbprint-list 6938fd4d98bab03faadb97b34396831e3780aea1 \
    --region ap-south-1

echo "OIDC Provider created successfully!"

# Create trust policy for the role
cat > trust-policy.json << EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Federated": "arn:aws:iam::172149903700:oidc-provider/token.actions.githubusercontent.com"
            },
            "Action": "sts:AssumeRoleWithWebIdentity",
            "Condition": {
                "StringEquals": {
                    "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
                },
                "StringLike": {
                    "token.actions.githubusercontent.com:sub": "repo:sauravsinghs/simplebank:*"
                }
            }
        }
    ]
}
EOF

# Update the role trust policy
aws iam update-assume-role-policy \
    --role-name github-ci-deploy \
    --policy-document file://trust-policy.json \
    --region ap-south-1

echo "Trust policy updated for github-ci-deploy role!"

# Clean up
rm trust-policy.json

echo "OIDC setup complete! You can now switch back to role-to-assume in the workflow." 