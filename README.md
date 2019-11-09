# link-of-the-day

## What is link-of-the-day?

Link of the day, well, does what it says on the tin. It redirects a user to a different link every day. I built this for [QuantumBlack](https://quantumblack.com), with the plan to generate a QR code that encodes `qb.zainp.com` (where my version of this is hosted) and build the QR code with post-it notes. People can walk by and scan the QR code and get a different link everyday, where links are crowdsources from employees. We'll see if it catches on...

## Why couldn't you come up with a better name?

I'm a software engineer and a mathematician, that should answer your question...

## How do I run it?

`git clone` the repository, make sure you have Docker :whale: installed, then `docker-compose up` should do it.

## How does it work?

It's a simple Go HTTP server, connected to a Postgres database and backed by the `links` package that retrieves the URL scheduled for the current date, updates the visit count and saves new URLs. I deploy it on [Heroku](https://heroku.com) and manage the domain via [Netlify](https://netlify.com).

A `GET` request will be met with a `301` temporary redirect, to the days scheduled link and a `POST` or `PUT` request, appropriate formatted with `Content-Type application/x-www-form-urlencoded` with data `url=<YOUR-URL>` will save the URL and automatically schedule it for the next available slot, so it's very much first come, first serve.