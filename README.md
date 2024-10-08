# Magic Init 

<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/2/25/Wizard_hat_and_wand.svg/1000px-Wizard_hat_and_wand.svg.png?20230507221104" height="200">

MagicInit is a Dagger Module that runs inference on your project and generates a relevant starter module. Go from 0 to pipeline with one command. 

It's not magic. It's Dagger.™

# Usage 

Infer project and generate dagger pipeline along with a GitHub Actions workflow that will execute your Dagger pipeline.

```
dagger call magicinit --source . --output .
```

# Sample Projects for Testing 

Go -> https://github.com/DiceDB/dice 
Typescript -> https://github.com/medplum/medplum
Python -> https://github.com/Infisical/python-sdk-official
Ruby -> https://github.com/mastodon/mastodon

# Future Improvements 
* be able to create the dir ourselves in the future 
* be able to specify which SDK you want to use 
* make Dagger version in template dagger.json be more intelligent
* add this to dagger init by default 
    * will allow us to not have to include --source
    * will allow to add dagger.json to the repo root
