import * as React from "react";
import { createRoot } from "react-dom/client";
import { Header } from "./header.js";
import { Listings } from "./listings.js";
import { Posting } from "./posting.js";

const URL = "http://localhost:5000";
const root_div = document.getElementById("root");
const root = createRoot(root_div);

function Main() {
  const [showListing, updateListing] = React.useState(true);

  return (
    <div>
      <Header updateListing={updateListing} />
      {showListing && <Listings url={URL} />}
      {!showListing && <Posting url={URL} />}
    </div>
  );
}

root.render(
  <div>
    <Main />
  </div>,
);
