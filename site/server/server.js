const express = require("express");
const cors = require("cors");
const multer = require("multer");

const app = express();
const storage = multer.memoryStorage()
const mult = multer({ storage: storage });
app.use(cors());

app.get("/total_listings", (req, res) => {
  res.json({ total_listings: 20 });
});

app.get("/listings/:page", (req, res) => {
  const apt_listings = [
    {
      _id: 1,
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
      _id: 2,
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
      _id: 3,
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
      _id: 4,
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

app.post("/new_listing", mult.single("apt_img"), (req, res) => {
  console.log(req.file);
  console.log("received listing!\n");
  console.log(req.body);
  res.status(201).send();
});

app.listen(5000, () => console.log("server is up and listening on port 5000"));
