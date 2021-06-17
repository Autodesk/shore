# Shore Cleanup Usage Doc

## What is a cleanup pipeline?

When running pipelines, there are *usually* resources created.

Those resources may be (examples):

- VM's (I.E. EC2)
- ECS tasks
- K8S Deployments
- Database Entries (I.E. new rows in a DynamoDB table)
- Lambda functions
- Managed Infra:
  - SQS
  - SNS
  - Databases
  - Certificates

And countless other resources.

When creating these resources, we tend to neglect the cleanup process of these resources.

## Shore Cleanup

The cleanup functionality requires developers to consider the opposite (reverse) operation of creating resources.

This is done by requiring the existence of a file name `cleanup.pipeline.jsonnet`.

As the name suggest, this is where the `cleanup` logic is implemented.

### Features

`shore cleanup` follows the same structure and logic of `shore render/save/exec`.

This design allows us to minimize context switches between running the `main.pipeline` & `cleanup.pipeline`.

#### TODO

`cleanup` logic can be optionally run after E2E tests (`shore test-remote`).

Allowing for a clean teardown logic.
