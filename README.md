# Sting of the Viper

This demonstrates how to integrate spf13/cobra with spf13/viper such command-line flags have the highest precedence,
then environment variables, then config file values, and then defaults set on command-line flags.

📬 Read the accompanying blog post that explains this example code. [Sting of the Viper: Getting Cobra and Viper to work together](https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/)

It also handles binding command-line flags that have dashes properly to environment variables with underscores. 
For example, `--favorite-color` is set with the environment variable `FAVORITE_COLOR`.

# Try it out

First grab the source code with go get or by cloning it. Change into the directory of this repository.

```
go get github.com/carolynvs/stingoftheviper
# or
git clone https://github.com/carolynvs/stingoftheviper.git
cd stingoftheviper/
```

Now let's build the CLI (stingoftheviper) and make sure everything is still working:

```
go build .
go test ./...
```

We are now ready to try out a few scenarios to test out the precedence order. First let's run it with no flags
or environment variables.

```console
$ ./stingoftheviper
Your favorite color is: blue
The magic number is: 7
```

If you take a peek at the config file, you will see that only the favorite-color was set there. So favorite-color
got its value from the config file, while magic number got its value from the flag's default value set in main.go.
So the lowest precedence is the flag default, followed by the config file.

Let's try setting an environment variable.

```console
$ STING_FAVORITE_COLOR=purple ./stingoftheviper
Your favorite color is: purple
The magic number is: 7
```

There's two interesting things going on here. One is that the environment variable has higher precedence than the
config file value obviously. The other is that not only was the environment variable automatically bound to
the flag, but we handled swapping the dashes for underscores in the binding (which isn't done for us in the 
library, you have to do that yourself).

To finish things off, let's actually use a flag.

```console
$ ./stingoftheviper --number 2
Your favorite color is: blue
The magic number is: 2
```

MAGIC! 🎩✨
