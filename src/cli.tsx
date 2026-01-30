#!/usr/bin/env node
import React from "react";
import { render } from "ink";
import meow from "meow";
import { App } from "./app.js";

const cli = meow(
  `
  Usage
    $ gaw

  Options
    --interval, -i  Polling interval in seconds (default: 10)

  Examples
    $ gaw
    $ gaw --interval 5
    $ gaw -i 30
`,
  {
    importMeta: import.meta,
    flags: {
      interval: {
        type: "number",
        shortFlag: "i",
        default: 10,
      },
    },
  },
);

render(<App interval={cli.flags.interval * 1000} />);
