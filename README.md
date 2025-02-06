# Triforce Blitz

This repository houses the Triforce Blitz web generator and its related 
suite of command line tools. If you're looking for the Ocarina of Time
Randomizer mod by Elagatua, please visit 
[this repository](https://github.com/Elagatua/OoT-Randomizer).

This project is hosting the current, ongoing rewrite of the website. The
old web generator was made in less than two days and the source code was 
lost. Despite that, it is still running to this day without any issue
over [triforceblitz.com](https://www.triforceblitz.com). The original
generator was written in Java using Spring Boot; the rewrite is written in
[Go](https://go.dev).

This project is directly affiliated with Triforce Blitz and is the official
web generator for the project.

## Building

Clone the repository to your computer and compile using `go`.

```bash
go mod download

# Manually build the executables and install them into $GOBIN (recommended)
go build -o $(go env GOBIN)/triforceblitz-server ./cmd/server
go build -o $(go env GOBIN)/triforceblitz-updater ./cmd/updater

# Alternatively, install all at once with default names (not recommended)
go install ./...
```

## Usage

For advanced usage details please see 
[the Wiki](https://github.com/jessebrands/triforceblitz/wiki/Usage).

The following example installs all generators and runs the web application
server on the default address.

```bash
# Installs all generators from the 'blitz' branch.
triforceblitz-updater install

# Run the application server.
triforceblitz-server
```

Open your browser and go to [localhost:8000](http://localhost:8000), if
everything went correctly, you should see the Triforce Blitz website.

## Frequently Asked Questions

### What is this project?

This project is, confusingly, not the Triforce Blitz project, but rather
the official web generator for the Triforce Blitz project. 

### What is the history of this project?

Generating seeds for Triforce Blitz by hand used to be slightly cumbersome 
and annoying, and in the early days, one had to manually download the 
generator, install Python on their computers, and generate seeds using the 
slightly unintuitive interface of the randomizer. This raised the need for
a convenient web generator which would make it easy to try Triforce Blitz
as well as sharing seeds with other players to race them competitively.

The author of the web generator was, at the time, working on a 
JavaScript-based ROM patcher for N64 Zelda ROMs called `zelda64.js`. When he
learned of the Triforce Blitz project from Elagatua, he figured that creating
a simple web generator for Triforce Blitz would be a great opportunity to try
out _zelda64.js_ in a production capacity. The generator quickly became very
popular and soon after became the de facto way for people to generate and
share seeds, though _zelda64.js_ has since been discontinued.

These days, the web generator remains very popular and the primary way people
roll seeds for Triforce Blitz, despite the fact that the generator is now also
available through the
[official Ocarina of Time Randomizer web generator](https://ootrandomizer.com),
thanks to it offering several features the community really enjoys such as 
automatic spoiler unlocks through 
[Racetime.gg](https://racetime.gg) integration, as well as the ever popular,
automatically generated _Seed of the Day_.

### What is Triforce Blitz?

[Triforce Blitz](https://www.triforceblitz.com) is an exciting, fast-paced,
and competitive take on 
[Ocarina of Time Randomizer](https://www.ootrandomizer.com). It is a set of
modifications to the original randomizer that seeks to create a ruleset that
is meant to provide a shorter gameplay session with a higher focus on solving
the generated game through logical deduction. It is created and maintained by
[Elagatua](https://github.com/Elagatua).

The official Discord community for Triforce Blitz can be found at the 
[ZeldaSpeedRuns Discord server](https://discord.gg/pZx9cpM7D2), you can find
us there in the `#tfb-general` channel. The official Triforce Blitz tournament
is hosted by
[The Silver Gauntlets](https://discord.gg/s5Bd23xeX9).

### What is Ocarina of Time Randomizer?

[Ocarina of Time Randomizer](https://www.ootrandomizer.com) is a mod for 
[The Legend of Zelda: Ocarina of Time](https://en.wikipedia.org/wiki/The_Legend_of_Zelda:_Ocarina_of_Time)
that, like many other _"item randomizers"_, modifies the game data by 
changing the location of items in the game as well as other aspects of 
the game such as starting location. To make sure the experience isn't
frustrating, a logic solver is used by the randomizer to ensure that 
every seed is completable without the use of glitches. This essentially
turns the game into a puzzle to be solved that is unique for every
generated seed.

There is a large, welcoming community of players over at the
[Ocarina of Time Randomizer Discord](https://discord.gg/ootrandomizer),
furthermore live races are hosted at
[Racetime.gg](https://racetime.gg) where players compete with one another
to beat the seed as quickly as possible. You can usually find  people 
streaming their gameplay (and matches) over on
[Twitch](https://www.twitch.tv/directory/category/the-legend-of-zelda-ocarina-of-time?tl=Randomizer).
Lastly, competitive tournaments are hosted regularly by Zelda speed running
communities such as
[ZeldaSpeedRuns](https://www.zeldaspeedruns.com) and
[The Silver Gauntlets](https://www.twitch.tv/thesilvergauntlets/about).

## Contributing

Pull requests are welcomed. For major changes, please open a GitHub issue
or contact `@0x0BEE` on Discord to discuss what you would like to change.

Please ensure that you write tests (where possible) or your PR will most
likely be rejected. Before submitting your pull request, ensure that all
tests pass and that you can successfully build the Dockerfile.

## License

This project is licensed under the 
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.html).
