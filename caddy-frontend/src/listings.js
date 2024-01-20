import * as React from "react";
import Pagination from "@mui/material/Pagination";
import Typography from "@mui/material/Typography";
import Stack from "@mui/material/Stack";
import Card from "@mui/material/Card";
import CardActions from "@mui/material/CardActions";
import CardMedia from "@mui/material/CardMedia";
import CardContent from "@mui/material/CardContent";
import Button from "@mui/material/Button";
import Box from "@mui/material/Box";
import Modal from "@mui/material/Modal";
import rconfidence from "./rent-with-confidence.jpg";

const boxModalStyle = {
  position: "absolute",
  top: "50%",
  left: "50%",
  transform: "translate(-50%, -50%)",
  width: 400,
  bgcolor: "background.paper",
  border: "2px solid #000",
  boxShadow: 24,
  p: 4,
};

function DisplayCard() {
  return (
    <div
      style={{
        backgroundPosition: "center",
        backgroundSize: "cover",
        backgroundImage: `url(${rconfidence})`,
      }}
      className="flex justify-center items-center w-full h-64 border-b-2 border-black"
    >
      <h1 className="text-8xl text-white">Rent with Confidence!</h1>
    </div>
  );
}

function BasicModal(props) {
  const [open, setOpen] = React.useState(false);
  const handleOpen = () => setOpen(true);
  const handleClose = () => setOpen(false);

  return (
    <div>
      <Button onClick={handleOpen}> More info </Button>
      <Modal
        open={open}
        onClose={handleClose}
        aria-labelledby={"modal-modal-title" + props.listing.id_}
        aria-describedby={"modal-modal-description" + props.listing.id_}
      >
        <Box sx={boxModalStyle}>
          <Typography
            id={"modal-modal-title" + props.listing.id_}
            variant="h6"
            component="h2"
          >
            {props.listing.listing_name}
          </Typography>
          <Typography
            id={"modal-modal-description" + props.listing.id_}
            sx={{ mt: 2 }}
          >
            {props.listing.listing_name} (Apartment number:{" "}
            {props.listing.address1}) is located at {props.listing.address2},{" "}
            {props.listing.pincode}. <br />
            Contact {props.listing.owner_email} from more details.
          </Typography>
        </Box>
      </Modal>
    </div>
  );
}

function DisplayListings(props) {
  const listings = props.listings.map((listing) => {
    return (
      <div key={listing._id} className="my-5">
        <Card sx={{ maxWidth: 400 }}>
          <CardMedia
            component="img"
            alt="listing image"
            height="140"
            image={rconfidence}
          />
          <CardContent>
            <Typography gutterBottom variant="h5" component="div">
              {listing.listing_name}
            </Typography>
            <Typography variant="h6" color="text.secondary">
              &#128176;| {listing.price}
            </Typography>
            <Typography variant="h6" color="text.secondary">
              &#128719;| {listing.rooms}
            </Typography>
          </CardContent>
          <CardActions>
            <BasicModal listing={listing} />
          </CardActions>
        </Card>
      </div>
    );
  });
  return (
    <div className="flex flex-row flex-wrap p-10 space-x-5">{listings}</div>
  );
}

async function getPages(url) {
  let pages = 0;
  try {
    const response = await fetch(`${url}/total_listings`);
    const tot_listings = (await response.json())["total_listings"];
    pages = tot_listings / 10 + (tot_listings % 10 != 0 ? 1 : 0);
  } catch (e) {
    console.error(e);
  }
  return pages;
}

async function getListings(url, page) {
  let listings = undefined;
  //console.log(page)
  try {
    const response = await fetch(`${url}/listings/${page - 1}`);
    listings = (await response.json())["listings"];
  } catch (e) {
    console.error(e);
  }
  return listings;
}

export function Listings(props) {
  const [total_pages, updateTotalPages] = React.useState(1);
  const [listings, updateListings] = React.useState([]);
  const [page, updatePage] = React.useState(1);

  const handlePaginationChange = (event, value) => {
    updatePage(value);
  };
  const setPagination = () => {
    getPages(props.url).then((res) => {
      if (res > 0) updateTotalPages(res);
    });
  };
  const setListings = () => {
    getListings(props.url, page).then((res) => {
      if (listings.length > 0) updateListings(res);
    });
  };

  React.useEffect(setPagination, []);
  React.useEffect(setListings, [page]);

  return (
    <div className="w-full max-w-full typo-bg min-h-screen">
      <DisplayCard />
      <DisplayListings listings={listings} />
      <Stack className="fixed top-90 left-1/2 -translate-x-1/2" spacing={2}>
        <Pagination
          count={total_pages}
          shape="rounded"
          variant="outlined"
          onChange={handlePaginationChange}
        />
      </Stack>
    </div>
  );
}
