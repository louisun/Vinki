import React, { Component } from 'react';

import classNames from 'classnames';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import {
  IconButton,
  Snackbar,
  SnackbarContent,
  withStyles,
} from '@material-ui/core';
import CheckCircleIcon from '@material-ui/icons/CheckCircle';
import CloseIcon from '@material-ui/icons/Close';
import ErrorIcon from '@material-ui/icons/Error';
import InfoIcon from '@material-ui/icons/Info';
import WarningIcon from '@material-ui/icons/Warning';

const mapStateToProps = state => {
  return {
    snackbar: state.snackbar,
  }
}

const mapDispatchToProps = () => { return {} }
const variantIcon = {
  success: CheckCircleIcon,
  warning: WarningIcon,
  error: ErrorIcon,
  info: InfoIcon,
};

const styles1 = theme => ({
  success: {
    backgroundColor: "#58BA82",
  },
  error: {
    backgroundColor: "#CC573A",
  },
  info: {
    backgroundColor: theme.palette.primary.dark,
  },
  warning: {
    backgroundColor: "#81A1C1",
  },
  icon: {
    fontSize: 20,
  },
  iconVariant: {
    opacity: 0.9,
    marginRight: theme.spacing(1),
  },
  message: {
    display: 'flex',
    alignItems: 'center',
    fontSize: "1.1em"
  },
})

function MySnackbarContent(props) {
  const { classes, className, message, onClose, variant, ...other } = props;
  const Icon = variantIcon[variant];

  return (
    <SnackbarContent
      className={classNames(classes[variant], className)}
      aria-describedby="client-snackbar"
      message={
        <span id="client-snackbar" className={classes.message}>
          <Icon className={classNames(classes.icon, classes.iconVariant)} />
          {message}
        </span>
      }
      action={[
        <IconButton
          key="close"
          aria-label="Close"
          color="inherit"
          className={classes.close}
          onClick={onClose}
        >
          <CloseIcon className={classes.icon} />
        </IconButton>,
      ]}
      {...other}
    />
  );
}
MySnackbarContent.propTypes = {
  classes: PropTypes.object.isRequired,
  className: PropTypes.string,
  message: PropTypes.node,
  onClose: PropTypes.func,
  variant: PropTypes.oneOf(['success', 'warning', 'error', 'info']).isRequired,
};

const MySnackbarContentWrapper = withStyles(styles1)(MySnackbarContent);
const styles = theme => ({
  margin: {
    margin: theme.spacing(1),
  },
})
class SnackbarCompoment extends Component {

  state = {
    open: false,
  }

  componentWillReceiveProps = (nextProps) => {
    if (nextProps.snackbar.toggle !== this.props.snackbar.toggle) {
      this.setState({ open: true });
    }
  }

  handleClose = () => {
    this.setState({ open: false });
  }

  render() {

    return (
      <Snackbar
        anchorOrigin={{
          vertical: this.props.snackbar.vertical,
          horizontal: this.props.snackbar.horizontal,
        }}
        open={this.state.open}
        autoHideDuration={6000}
        onClose={this.handleClose}
      >
        <MySnackbarContentWrapper
          onClose={this.handleClose}
          variant={this.props.snackbar.color}
          message={this.props.snackbar.msg}
        />
      </Snackbar>
    );
  }

}

const MessageBar = connect(
  mapStateToProps,
  mapDispatchToProps
)(withStyles(styles)(SnackbarCompoment))

export default MessageBar