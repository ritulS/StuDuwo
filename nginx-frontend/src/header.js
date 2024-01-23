import * as React from "react";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";

export function Header(props) {
  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static" color="primary">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            FakeRent.com
          </Typography>
          <Box>
            <Button color="inherit" onClick={() => props.updateListing(true)}>
              Listings
            </Button>
            <Button color="inherit" onClick={() => props.updateListing(false)}>
              Posting
            </Button>
          </Box>
        </Toolbar>
      </AppBar>
    </Box>
  );
}
