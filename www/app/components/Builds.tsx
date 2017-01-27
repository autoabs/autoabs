/// <reference path="../References.d.ts"/>
import * as React from 'react';
import AppBar from 'material-ui/AppBar';
import * as BuildTypes from '../types/BuildTypes';
import BuildStore from '../stores/BuildStore';
import * as BuildActions from '../actions/BuildActions';
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
		width: '100%',
		flex: 1,
		display: 'flex',
		flexDirection: 'row',
		flexWrap: 'wrap',
		padding: '5px',
	} as React.CSSProperties,
};

export default class Builds extends React.Component<null, State> {
	syncing: boolean;

	constructor(props: any, context: any) {
		super(props, context);
		BuildActions.sync();
		this.state = getState();
	}

	componentDidMount(): void {
		BuildStore.addChangeListener(this.onChange);
		this.syncing = true;
		this.sync();
	}

	componentWillUnmount(): void {
		BuildStore.removeChangeListener(this.onChange);
		this.syncing = false;
	}

	onChange = (): void => {
		this.setState(getState());
	}

	sync = (): void => {
		setTimeout(() => {
			if (!this.syncing) {
				return;
			}

			BuildActions.sync().then(
				this.sync,
				this.sync,
			);
		}, 1000);
	}

	render(): JSX.Element {
		let builds = this.state.builds;

		let buildsDom: JSX.Element[] = [];
		for (let build of builds) {
			buildsDom.push(<Build key={build.id} build={build}/>);
		}

		return <div>
			<AppBar title="AutoABS"/>
			<div style={css.builds}>
				{buildsDom}
			</div>
		</div>;
	}
}
