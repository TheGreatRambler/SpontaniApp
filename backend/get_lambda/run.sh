aws lambda create-function --function-name getLambda \
    --runtime provided.al2023 --handler bootstrap \
    --architectures arm64 \
    --role arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole \
    --zip-file fileb://getLambda.zip \
    --region us-east-2