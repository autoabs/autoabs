/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as BuildTypes from '../types/BuildTypes'
import BuildStore from '../stores/BuildStore';
import * as BuildUtils from '../utils/BuildUtils';
import Build from './Build';

interface State {
	builds: BuildTypes.Builds;
}

function getState(): State {
	return {
		builds: BuildStore.builds,
	};
}

const css = {
	builds: {
		display: 'flex',
		flexDirection: 'row',
		flexWrap: 'wrap',
		padding: '5px',
	} as React.CSSProperties,
};

export default class Builds extends React.Component<null, State> {
	constructor(props: any, context: any) {
		super(props, context);
		BuildUtils.init();
		this.state = getState();
	}

	componentDidMount(): void {
		BuildStore.addChangeListener(this._onChange);
	}

	componentWillUnmount(): void {
		BuildStore.removeChangeListener(this._onChange);
	}

	_onChange = (): void => {
		this.setState(getState());
	};

	render(): JSX.Element {
		let builds = this.state.builds;

		let buildsDom: JSX.Element[] = [];
		for (let key in builds) {
			buildsDom.push(<Build key={key} build={builds[key]}/>);
		}

		return <div style={css.builds}>
			{buildsDom}
		</div>;
	}
}
