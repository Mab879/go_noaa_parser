# go_noaa_praser

This program parsers the NOAA Warning Atom Feed and saves the data to a database.


## How to run
1. Clone the repository
2. Run `go get` to install the dependencies
3. Run `go build` to build the program
4. Create the database PostgreSQL database
5. Migrate the database by running `psql -d <database> < create_table.sql`
6. Set the environment variable `DATABASE_URL` to the database URL
   1. Example: `export DATABASE_URL=postgres://user:password@localhost:5432/database`
7. Run the program `./go_noaa_parser`


## LICENSE
This project is licensed under the GNU General Public License, version 2.0.
See the LICENSE file for more details.