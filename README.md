# Gohm ॐ ![build](https://travis-ci.org/pote/gohm.svg)

Gohm is a Go port of the popular [Ohm](https://github.com/soveran/ohm) Ruby library, it provides a simple interface to store and retrieve your model data in a Redis database.

## Ohm compatibility

Gohm implements nothing but the basic usage right now, but expect all or most features in Ohm to be implemented into Gohm as time goes by, contributions are very welcome. :)

Both Ohm and Gohm are powered by [ohm-scripts](https://github.com/soveran/ohm-scripts), a set of Lua scripts that bundle common operations and make it easy to write a port such as this one, it also means that by adhering to the ohm standard **models stored with Gohm can be loaded from Ohm, and vice-versa**.

## Gohm Models

Gohm models are simple Go structs which just need to adhere certain criteria, this is how a Gohm model might look like in code:

```go
type user struct{
	ID    string `ohm:"id"`
	Name  string `ohm:"name"`
	Email string `ohm:"email"`
}
```

As you can see, when you describe your Go struct you specify the fields that Gohm will be persisting by identifying them with tags, just as you would describe JSON serialization, for example.

The two special considerations you need to keep in your struct for it to be Gohm-compliant are:

1) It must have a field tagged as `id`.
2) All fields tagged with `ohm` need to be public (capitalized)\*

* Only necessary because of how Go's reflection works, Gohms need to be able to manipulate the attribute and it's value.


### Setting up the Gohm connection.

You can set up a new connection easily:

```go
import(
  "github.com/pote/gohm"
)

// Set up using the `REDIS_URL` env variable.
Gohm, _ := gohm.NewConnection()

// Or use your own Redis pool.
import(
  "github.com/pote/gohm"
  "github.com/pote/redisurl"
)

pool := redisurl.NewPoolWithUrl("redis://:password@localhost:6379", 3, 200, "240s")
Gohm, _ := gohm.NewConnection(pool)
```

### Saving a model

```go
u := User{
  Name: "Marty",
  Email: "marty@mcfly.com",
}

Gohm.Save(&u)

u.ID //=> "1"
```

### Retrieving a model by id

```go
u := User{ID: "1"}
Gohm.Load(&u)

u.Name //=> "Marty"
```
