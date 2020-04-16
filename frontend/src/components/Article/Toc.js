import React from 'react';

import { makeStyles } from '@material-ui/core';
import Typography from '@material-ui/core/Typography';

// 实际是 css
const useStyles = makeStyles((theme) => ({
  right: {
    flex: "0 0 250px",
  },
  rightWrapper: {
    position: "fixed",
    top: "50px",
    right: "0",
  },
  toc: {
    marginTop: "30px",
    padding: "0 20px 0 20px",
    maxWidth: "250px",
    color: "#4C566A",
  },
}));

export default function Toc() {
  const classes = useStyles()  
  return (
    <div className={classes.right}>
      <div className={classes.rightWrapper}>
        <div className={classes.toc}>
          <Typography variant="subtitle1" gutterBottom>
            subtitle1. Lorem ipsum dolor sit amet, consectetur adipisicing elit.
            Quos blanditiis tenetur
          </Typography>
          <Typography variant="subtitle2" gutterBottom>
            subtitle2. Lorem ipsum dolor sit amet, consectetur adipisicing elit.
            Quos blanditiis tenetur
          </Typography>
          <Typography variant="body1" gutterBottom>
            body1. Lorem ipsum dolor sit amet, consectetur adipisicing elit.
            Quos blanditiis tenetur unde suscipit, quam beatae rerum inventore
            consectetur, neque doloribus, cupiditate numquam dignissimos laborum
            fugiat deleniti? Eum quasi quidem quibusdam.
          </Typography>
          <Typography variant="body2" gutterBottom>
            body2. Lorem ipsum dolor sit amet, consectetur adipisicing elit.
            Quos blanditiis tenetur unde suscipit, quam beatae rerum inventore
            consectetur, neque doloribus, cupiditate numquam dignissimos laborum
            fugiat deleniti? Eum quasi quidem quibusdam.
          </Typography>
        </div>
      </div>
    </div>
  );
}
