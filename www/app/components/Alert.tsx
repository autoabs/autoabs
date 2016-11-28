/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Snackbar from 'material-ui/Snackbar';
import * as AlertTypes from '../types/AlertTypes';
import * as AlertActions from '../actions/AlertActions';

interface Props {
	alert: AlertTypes.Alert;
}

const css = {
	bar: {
		position: 'static',
		left: 'auto',
		bottom: 'auto',
		zIndex: 'auto',
		marginTop: '5px',
	} as React.CSSProperties,
	body: {
		height: 'auto',
	} as React.CSSProperties,
};

export default class Alert extends React.Component<Props, void> {
	onTouch = (): void => {
		AlertActions.remove(this.props.alert.id);
	}

	render(): JSX.Element {
		return <Snackbar
			style={css.bar}
			bodyStyle={css.body}
			open={true}
			action="Close"
			message={this.props.alert.message}
			onActionTouchTap={this.onTouch}
			onRequestClose={() => null}
		/>;
	}
}
