# react-alert-template-mui

## Installation

```
npm install react-alert-template-mui
```

## Demo

[![Edit react-alert-template-mui example](https://codesandbox.io/static/img/play-codesandbox.svg)](https://codesandbox.io/s/react-alert-template-mui-example-up2ys?fontsize=14)

## How to Use

### Simple alert

```javascript
alert.show("Oh look, an alert!");
```

### Alert with title

```javascript
alert.show("This is an alert with title!", {
  title: "Random Alert Title",
});
```

### Change copy on close button

```javascript
alert.show("This is an alert with title!", {
  title: "Random Alert Title",
  closeCopy: "Cancel",
});
```

### Alert with extra actions

```javascript
alert.show("This is an alert with extra actions!", {
  title: "Alert with extra actions!",
  actions: [
    {
      copy: "Do something",
      onClick: () => setAction("Actioned!"),
    },
  ],
});
```
