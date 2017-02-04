/// <reference path="../References.d.ts"/>
import * as React from 'react';
import AppBar from 'material-ui/AppBar';
import * as BuildTypes from '../types/BuildTypes';
import BuildStore from '../stores/BuildStore';
import InfiniteFlex from './InfiniteFlex';
import * as BuildActions from '../actions/BuildActions';
import Build from './Build';

interface State {
	builds: BuildTypes.Builds;
	count: number;
}

function getState(): State {
	return {
		builds: BuildStore.builds,
		count: BuildStore.count,
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
		this.state = getState();
		BuildActions.sync();
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

	buildItem = (index: number, build: BuildTypes.Build): JSX.Element => {
		return <Build key={index} build={build}/>
	}

	render(): JSX.Element {
		return <div>
			<AppBar title="AutoABS"/>
			<InfiniteFlex
				style={css.builds}
				width={260}
				height={123}
				padding={5}
				scrollMargin={2}
				scrollMarginHit={1}
				buildItem={this.buildItem}
				items={this.state.builds}
				count={this.state.count}
			/>
		</div>;
	}
}
