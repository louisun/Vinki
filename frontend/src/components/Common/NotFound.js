import React from 'react';

import {
  lighten,
  makeStyles,
} from '@material-ui/core/styles';
import SentimentVeryDissatisfiedIcon
  from '@material-ui/icons/SentimentVeryDissatisfied';

const useStyles = makeStyles(theme=>({
    icon:{
        fontSize: "160px",
    },
    emptyContainer: {
        bottom: "0",
        height: "500",
        margin: "auto",
        width: "500",
        color: lighten(theme.palette.text.disabled,0.4),
        textAlign: "center",
        paddingTop: "20px"
    },
    emptyInfoBig: {
        fontSize: "35px",
        color: lighten(theme.palette.text.disabled,0.4),
    },
}));

export default function NotFound(props) {
    const classes = useStyles();
    return (
        <div className={classes.emptyContainer}>
            <SentimentVeryDissatisfiedIcon className={classes.icon}/>
            <div className={classes.emptyInfoBig}>
                {props.msg}
            </div>
        </div>

    )
}