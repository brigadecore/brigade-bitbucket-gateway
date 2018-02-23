# Scripting Guide

This guide explains the basics of events available to `brigade.js` files.

For more, see the [Brigade Scripting Guide](https://github.com/Azure/brigade/blob/master/docs/topics/scripting.md)

# Brigade Bitbucket Events

Brigade listens for certain things to happen, this gateway provides those such events from a Bitbucket repository. The events that Brigade listens for are configured in your project.

When Brigade observes such an event, it will load the `brigade.js` file and see if there is an event handler that matches the event.

For example:

```javascript
const { events } = require("brigadier")

events.on("push", () => {
  console.log("==> handling an 'push' event")
})
```

The Bitbucket Gateway produces 6 events:

```
push
repo:fork
repo:updated
repo:commit_comment_created
repo:commit_status_created
repo:commit_status_updated

```

These are based on the events described in the [Bitbucket Webhooks API](https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html) guide.