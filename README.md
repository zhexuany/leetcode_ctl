 Leetcode-ctl
 ===========
 
 Leetcode-ctl is a command line controller for submitting problem solution to leetcode online judge. The goal of this project is to provide a simple way
 of grabing leetcode problem and submiting your solution just in terminal. 

# Build #
A simple commnad `make` can do the building.

# How to use this project #
Before you using this project, you have to grab cookie from leetcode website. You can google about how to get use data.
A simple config looks like the following:
~~~ 
leetcode-session= "your leetcode-session"
csrf-token= "your csrf-token"
lang-type= "your language"
~~~

Recently, I use `go` a lot. So my config's `lang-type` is `golang`. You can specify your faviorate language. 

A function named `isValidLanguageType` will valid your `lang-type` before all operation. 

~~~
func isValidLanguageType(language string) bool {
	switch language {
	case "golang", "java", "csharp", "cpp", "c", "javascript":
		return true
	}
	return false
}
~~~

## Generating a problem file ##

~~~
./leetcode-ctl generate -id 1 -config default.toml
~~~

## Submiting a solution ##

~~~
./leetcode-ctl submit -file two-sum -id 1 -config default.toml
~~~




