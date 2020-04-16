import React from 'react';

import Chip from '@material-ui/core/Chip';
import { makeStyles } from '@material-ui/core/styles';
import KeyboardArrowRightIcon from '@material-ui/icons/KeyboardArrowRight';

const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
    justifyContent: 'left',
    flexWrap: 'wrap',
    padding: theme.spacing(0.5),
  },
  chip: {
    margin: theme.spacing(1),
    backgroundColor: "#E5E9F0",
    color: "#4C566A",
    fontSize: "15px",
    // fontWeight: "bold",
  },
}));

export default function ChipsArray() {
  const classes = useStyles();
  const [chipData, setChipData] = React.useState([
    { key: 0, label: 'Angular' },
    { key: 1, label: 'jQuery' },
    { key: 2, label: 'Polymer' },
    { key: 3, label: 'React' },
    { key: 4, label: 'Vue.js' },
    { key: 5, label: 'Vue.js' },
    { key: 6, label: 'Vue.js' },
    { key: 7, label: 'Vue.js 从入门到退出' },
    { key: 8, label: 'Vue.js' },
    { key: 9, label: 'Vue.js' },
    { key: 10, label: 'Vue.js' },
    { key: 11, label: 'Vue.js' },
    { key: 12, label: 'Vue.js' },
    { key: 13, label: 'Vue.js' },
    { key: 14, label: 'Vue.js' },
    { key: 15, label: 'Vue.js' },
    { key: 16, label: 'Vue.js' },
  ]);

  return (
    <div className={classes.root}>
      {chipData.map((data) => {
        let icon = <KeyboardArrowRightIcon/>;

        return (
          <Chip
            key={data.key}
            icon={icon}
            label={data.label}
            className={classes.chip}
          />
        );
      })}
    </div>
  );
}