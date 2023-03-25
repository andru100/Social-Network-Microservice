import React, { Fragment } from "react";
import { useAlert } from "react-alert";

const ProcessErrorAlerts = () => {
  const alert = useAlert();


  return (
    <Fragment>
      <button
        onClick={() => {
          alert.show("Oh look, an alert!");
        }}
      >
        Show Alert
      </button>
      <button
        onClick={() => {
          alert.error("You just broke something!");
        }}
      >
        Oops, an error
      </button>
      <button
        onClick={() => {
          alert.success("It's ok now!");
        }}
      >
        Success!
      </button>
    </Fragment>
  );
};

export default ProcessErrorAlerts;