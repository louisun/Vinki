import React from 'react';

import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import IconButton from '@material-ui/core/IconButton';
import { fade, makeStyles } from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar';
import InputBase from '@material-ui/core/InputBase';
import Typography from '@material-ui/core/Typography';
import MenuIcon from '@material-ui/icons/Menu';
import MailIcon from '@material-ui/icons/Mail';
import SearchIcon from '@material-ui/icons/Search';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import InboxIcon from '@material-ui/icons/MoveToInbox';
import DraftsIcon from '@material-ui/icons/Drafts';
import SendIcon from '@material-ui/icons/Send';
import { withStyles } from '@material-ui/core/styles';
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';

const StyledMenu = withStyles({
  paper: {
    border: '1px solid #d3d4d5',
  },
})((props) => (
    <Menu
        elevation={0}
        getContentAnchorEl={null}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'center',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'center',
        }}
        {...props}
    />
));

const StyledMenuItem = withStyles((theme) => ({
  root: {
    '&:focus': {
      backgroundColor: theme.palette.primary.main,
      '& .MuiListItemIcon-root, & .MuiListItemText-primary': {
        color: theme.palette.common.white,
      },
    },
  },
}))(MenuItem);

const useStyles = makeStyles((theme) => ({
  root: {
      display: "flex",
  },
  appbar: {
    flexGrow: 1,
    height: "50px",
  },
  menuButton: {
    marginRight: theme.spacing(2),
  },
  brand: {
    color: "#86C0D2",
    fontWeight: "bold",
  },
  repo: {
    fontSize: "1.1rem",
    color: "#8FBCBB",
    fontWeight: "bold",
    background: '#4E5668',
    boxShadow: 'none',
    textTransform: 'none',
    '&:hover': {
      backgroundColor: '#81A1C1',
      color: '#E5E9F1',
      borderColor: '#4E5668',
      boxShadow: 'none',
    },
    '&:active': {
      color: '#E5E9F1',
      backgroundColor: '#81A1C1',
      borderColor: '#4E5668',
      boxShadow: 'none',
    },
  },
  grow: {
    flexGrow: 1,
  },
  search: {
    position: 'relative',
    borderRadius: theme.shape.borderRadius,
    backgroundColor: fade("#F2F4F8", 0.15),
    '&:hover': {
      backgroundColor: fade(theme.palette.common.white, 0.25),
    },
    marginLeft: 0,
    marginRight: "30px",
    width: '100%',
    [theme.breakpoints.up('sm')]: {
      marginLeft: theme.spacing(1),
      width: 'auto',
    },
    [theme.breakpoints.down('sm')]: {
        display: 'none',
    },
  },
  searchIcon: {
    padding: theme.spacing(0, 2),
    height: '100%',
    position: 'absolute',
    pointerEvents: 'none',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  inputRoot: {
    color: 'inherit',
  },
  inputInput: {
    padding: theme.spacing(1, 1, 1, 0),
    // vertical padding + font size from searchIcon
    paddingLeft: `calc(1em + ${theme.spacing(4)}px)`,
    transition: theme.transitions.create('width'),
    width: '100%',
    [theme.breakpoints.up('sm')]: {
      width: '12ch',
      '&:focus': {
        width: '20ch',
      },
    },
  },
}));



export default function DefaultAppBar() {
  const classes = useStyles();
  const [anchorEl, setAnchorEl] = React.useState(null);

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };
  return (
    <div className={classes.root}>
      <AppBar position="fixed" color="primary" className={classes.appbar}>
        <Toolbar variant="dense">
          <IconButton edge="start" className={classes.menuButton} color="inherit" aria-label="menu">
            <MenuIcon />
          </IconButton>
          <Typography variant="h4" className={classes.brand}>
            Vinki
          </Typography>


          {/*<Typography variant="h6" className={classes.repo}>*/}
          {/*    Code*/}
          {/*</Typography>*/}
          {/*<div>*/}
          {/*  <IconButton aria-label="show 17 new notifications" color="inherit">*/}
          {/*    <EmojiObjects />*/}
          {/*  </IconButton>*/}
          {/*</div>*/}
          <div className={classes.grow} />
          <div className={classes.search}>
            <div className={classes.searchIcon}>
              <SearchIcon />
            </div>
            <InputBase
                placeholder="搜索"
                classes={{
                  root: classes.inputRoot,
                  input: classes.inputInput,
                }}
                inputProps={{ 'aria-label': 'search' }}
            />
          </div>
          {/*<div>*/}
          {/*  <IconButton aria-label="show 4 new mails" color="inherit">*/}
          {/*      <MailIcon />*/}
          {/*  </IconButton>*/}
          {/*</div>*/}
          <div>
            <Button
                aria-controls="customized-menu"
                aria-haspopup="true"
                variant="contained"
                size="small"
                onClick={handleClick}
                endIcon={<KeyboardArrowDownIcon/>}
                className={classes.repo}
                disableRipple
            >
              Code
            </Button>
            <StyledMenu
                id="customized-menu"
                anchorEl={anchorEl}
                keepMounted
                open={Boolean(anchorEl)}
                onClose={handleClose}
            >
              <StyledMenuItem>
                <ListItemIcon>
                  <SendIcon fontSize="small" />
                </ListItemIcon>
                <ListItemText primary="Sent mail" />
              </StyledMenuItem>
              <StyledMenuItem>
                <ListItemIcon>
                  <DraftsIcon fontSize="small" />
                </ListItemIcon>
                <ListItemText primary="Drafts" />
              </StyledMenuItem>
              <StyledMenuItem>
                <ListItemIcon>
                  <InboxIcon fontSize="small" />
                </ListItemIcon>
                <ListItemText primary="Inbox" />
              </StyledMenuItem>
            </StyledMenu>
          </div>
        </Toolbar>
      </AppBar>
    </div>
  );
}