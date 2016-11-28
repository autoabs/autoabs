/// <reference path="../References.d.ts"/>
import * as React from 'react';
import AlertStore from '../stores/AlertStore';
import * as AlertTypes from '../types/AlertTypes';
import Alert from './Alert';

interface State {
	alerts: AlertTypes.Alerts;
}

function getState(): State {
	return {
		alerts: AlertStore.alerts,
	};
}

const css = {
	alerts: {
		position: 'fixed',
		left: '50%',
		bottom: 0,
		zIndex: 2900,
	} as React.CSSProperties,
};

export default class Alerts extends React.Component<void, State> {
	constructor(props: void, context: any) {
		super(props, context);
		this.state = getState();
	}

	componentDidMount(): void {
		AlertStore.addChangeListener(this.onChange);
	}

	componentWillUnmount(): void {
		AlertStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState(getState());
	}

	render(): JSX.Element {
		let alerts = this.state.alerts;

		let alertsDom: JSX.Element[] = [];
		for (let alert of alerts) {
			alertsDom.push(<Alert key={alert.id} alert={alert}/>);
		}

		return <div style={css.alerts}>
			{alertsDom}
		</div>;
	}
}
