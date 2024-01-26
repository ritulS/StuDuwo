import * as React from "react";
import { createRoot } from "react-dom/client";
import { Header } from "./header.jsx";
import { Listings } from "./listings.jsx";
import { Posting } from "./posting.jsx";

const URL = "https://api.348575.xyz";
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
