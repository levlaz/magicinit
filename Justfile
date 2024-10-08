# Helper file used for testing different projects
python:
    dagger call init --source https://github.com/Infisical/python-sdk-official --output .python

test-python:
    dagger call -m .python/.dagger --source https://github.com/Infisical/python-sdk-official ci 

typescript:
    dagger call init --source https://github.com/medplum/medplum --output .typescript

test-typescript:
    dagger -m .typescript/.dagger --source https://github.com/medplum/medplum call ci 

ruby:
    dagger call init --source https://github.com/aws/aws-sdk-ruby --output .ruby


test-ruby:
    dagger -m .ruby/.dagger --source https://github.com/aws/aws-sdk-ruby call ci 

go:
    dagger call init --source https://github.com/DiceDB/dice --output .go

test-go:
    dagger -m .go/.dagger --source https://github.com/DiceDB/dice call ci 

all:
    just python
    just typescript
    just ruby
    just go

clean:
    rm -rf .python
    rm -rf .typescript
    rm -rf .ruby
    rm -rf .go

edge:
    dagger -m github.com/TomChv/dagger@ts-sdk/rollback-gql-client call dev with-mounted-directory --path "/src" --source ".ruby/.dagger" with-workdir --path "/src" with-exec --args "dagger,call,--source,https://github.com/aws/aws-sdk-ruby,ci"