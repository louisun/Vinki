import React from 'react';

import { makeStyles } from '@material-ui/core';
import Avatar from '@material-ui/core/Avatar';
import Box from '@material-ui/core/Box';
import Button from '@material-ui/core/Button';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import CardMedia from '@material-ui/core/CardMedia';
import Divider from '@material-ui/core/Divider';
import Hidden from '@material-ui/core/Hidden';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemAvatar from '@material-ui/core/ListItemAvatar';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Typography from '@material-ui/core/Typography';
import FolderIcon from '@material-ui/icons/Folder';

import ChipsArray from './Chip';

// 实际是 css
const useStyles = makeStyles((theme) => ({
  leftWrapper: {
    position: "fixed",
    top: "50px",
    left: "0",
  },

  midWrapper: {
    margin: "auto",
    maxWidth: "1000px",
    minHeight: "70vh",
    padding: "20px 60px 20px 60px",
    backgroundColor: "#FFFFFF",
    color: "#4C566A",
    boxShadow:
      "0 4px 6px rgba(184,194,215,0.25), 0 5px 7px rgba(184,194,215,0.1)",
    borderRadius: "8px",
  },
  left: {
    flex: "0 0 250px",
  },
  mid: {
    marginTop: "30px",
    flex: "1 0",
  },

  card: {
    marginTop: "30px",
    maxWidth: "220px",
    marginLeft: "20px",
  },
  title: {
    verticalAlign: "middle",
    textIndent: "0.5em",
  },
  view: {
      marginTop: "30px",
      marginBottom: "30px",
  }
}));

function generate(element) {
  return [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15].map((value) =>
    React.cloneElement(element, {
      key: value,
    })
  );
}

export default function Article() {
  const classes = useStyles();

  return (
    <React.Fragment>
      <Hidden mdDown>
        <div className={classes.left}>
          <div className={classes.leftWrapper}>
            <Card className={classes.card}>
              <CardActionArea>
                <CardMedia
                  component="img"
                  alt="Contemplative Reptile"
                  height="140"
                  image="https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2020-04-14-00-33-02_r61.png"
                  title="Contemplative Reptile"
                />
                <CardContent>
                  <Typography gutterBottom variant="h5" component="h2">
                    Code
                  </Typography>
                  <Typography
                    variant="body2"
                    color="textSecondary"
                    component="p"
                  >
                    仓库
                  </Typography>
                </CardContent>
              </CardActionArea>
              <CardActions>
                <Button size="small" color="primary">
                  Share
                </Button>
                <Button size="small" color="primary">
                  Learn More
                </Button>
              </CardActions>
              <List>
                {generate(
                  <ListItem>
                    <ListItemAvatar>
                      <Avatar>
                        <FolderIcon />
                      </Avatar>
                    </ListItemAvatar>
                    <ListItemText primary="Golang" />
                  </ListItem>
                )}
              </List>
            </Card>
          </div>
        </div>
      </Hidden>

      <div className={classes.mid}>
        <div className={classes.midWrapper}>
          <Box className={classes.view}>
          <Typography variant="h5" gutterBottom className={classes.title}>
            子标签
          </Typography>
          <ChipsArray />
          </Box>
          <Divider variant="middle" />
          <Box className={classes.view}>
          <Typography variant="h5" gutterBottom className={classes.title}>
            文档列表
          </Typography>
          <div className={classes.demo}>
            <List>
              {generate(
                <ListItem>
                  <ListItemIcon>
                    <FolderIcon />
                  </ListItemIcon>
                  <ListItemText
                    primary="Single-line item"
                  />
                </ListItem>,
              )}
            </List>
          </div>
          </Box>
        </div>
      </div>
    </React.Fragment>
  );
}
