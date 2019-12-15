# stygis - golang Hexagonal Architecture 

Just another sample of hexagonal architecture for golang
focus on simplifed code, orgenized structure and better naming for functions and packages name
**No duplicate naming for package**

## Project structure

```sh
stygis/
├── bin                                 
├── cmd                                 # Packages that provide support for a specific program that is being built
      ├── rest                          # Entry point of main application based on REST technology
├── pkg                                 # Contains all application packages
      ├── dashboard                     # For declaring all constants variable based on it's entity, declaring with full name descripting it
            ├── query                   # where query is being declared for storage sql, and it's better using sqlx naming system for passing the parameters
            ├── state                   # contains all constants variable (string, int, etc..) like status, or constant strings
      ├── handler                       # any protocol interface that needs handler is being declare here, whether for cleaning input data, error payload handling, etc..
            ├── rest                    # handler for rest technology
      ├── internal                      # Domain packages are here, contains business logic and interfaces that belong to each domain, and no usecase calling other usecase
            ├── user                    # only sample domain, user package which handler for user business logic
                ├── initiator.go        # this is the only file to declare interface methods from storage and repository. also where to put func init the package.
                ....                    # name the file based on the usecase that may contains multiple acticity; example: login usecase for only one activity, profile contains edit profile and show profile
      ├── model                         # yes, don't write multiple structs, use tags and specific naming instead
      ├── repository                    # this may a bit weird for you, this package uses for data storing logic, it is an optional if your domain only saves data to one db, but it's different when a domain uses multiple storage, for example caching and multiple persistances
      ├── storage                       # this is where you put the data storing code. whether persistance like postgresql, monggodb, etc. and caching like redis, etc. 
            ├── cache                   # package where caching storing code is written based on its domain.
            ├── persistance             # like its name, this package contains storing code for SQL or noSQL db
```

## Notes

Be aware that this may not be the better structure for your application, but i'm hopping to write into more general structure and stay simple

## References
[Ashley McNamara + Brian Ketelsen. Go best practices.](https://youtu.be/MzTcsI6tn-0)