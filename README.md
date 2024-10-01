# Magic Init 

![](https://upload.wikimedia.org/wikipedia/commons/thumb/2/25/Wizard_hat_and_wand.svg/1000px-Wizard_hat_and_wand.svg.png?20230507221104)

MagicInit is a Dagger Module that runs inference on your project and generates a relevant starter module. Go from 0 to pipeline with one command. 

It's not magic. It's Dagger.™

# Usage 

Infer project and generate dagger pipeline.

```
dagger call magicinit -o .dagger
```

Optional args 

`--provider` default to gha, in the future be able to generate for others as well 
`--source` default to .

# Sample Projects for Testing 

Go -> https://github.com/DiceDB/dice 
Typescript -> https://github.com/medplum/medplum
Python -> https://github.com/Cinnamon/kotaemon
Ruby -> https://github.com/mastodon/mastodon

# Future Improvements 
* be able to create the dir ourselves in the future 
* be able to specify which SDK you want to use 
* make Dagger version in template dagger.json be more intelligent