const express = require("express");
const cors = require("cors");

const app = express();

app.use(cors());

app.get("/total_listings", (req, res) => {
  res.json({ total_listings: 20 });
});

app.get("/listings/:page", (req, res) => {
  const apt_listings = [
    {
      id_: 1,
      listing_name: "a",
      owner_email: "b",
      address1: "c",
      address2: "d",
      pincode: "e",
      apt_img: "f",
      price: "1000",
      rooms: "1",
    },
    {
      id_: 2,
      listing_name: "a2",
      owner_email: "b",
      address1: "c",
      address2: "d",
      pincode: "e",
      apt_img: "f",
      price: "1200",
      rooms: "2",
    },
    {
      id_: 3,
      listing_name: "a3",
      owner_email: "b",
      address1: "c",
      address2: "d",
      pincode: "e",
      apt_img: "f",
      price: "1300",
      rooms: "3",
    },
    {
      id_: 4,
      listing_name: "a4",
      owner_email: "b",
      address1: "c",
      address2: "d",
      pincode: "e",
      apt_img: "f",
      price: "1400",
      rooms: "4",
    },
  ];

  //console.log(req.params["page"])

  if (req.params["page"] == 0) {
    res.json({ listings: apt_listings.slice(0, 2) });
  } else {
    res.json({ listings: apt_listings.slice(2, 4) });
  }
});

app.post("/new_listing", (req, res) => {
  console.log("received listing!\n");
  console.log(req.body);
  res.status(201).send();
});

app.listen(5000, () => console.log("server is up and listening on port 5000"));
