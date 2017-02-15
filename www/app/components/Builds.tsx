/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as BuildTypes from '../types/BuildTypes';
import BuildStore from '../stores/BuildStore';
import InfiniteFlex from './InfiniteFlex';
import * as BuildActions from '../actions/BuildActions';
import BuildInfo from './BuildInfo';
import Loading from './Loading';
import Build from './Build';

interface State {
	builds: BuildTypes.Builds;
	index: number;
	count: number;
}

function getState(): State {
	return {
		builds: BuildStore.builds,
		index: BuildStore.index,
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
		this.syncLoop();
	}

	componentWillUnmount(): void {
		BuildStore.removeChangeListener(this.onChange);
		this.syncing = false;
	}

	onChange = (): void => {
		this.setState(getState());
	}

	syncLoop = (): void => {
		setTimeout(() => {
			if (!this.syncing) {
				return;
			}

			BuildActions.sync().then(
				this.syncLoop,
				this.syncLoop,
			);
		}, 1000);
	}

	sync = (): void => {
		BuildActions.sync();
	}

	buildItem = (index: number, build: BuildTypes.Build): JSX.Element => {
		return <Build key={index} build={build}/>
	}

	traverse = (index: number): void => {
		BuildActions.traverse(index);
	}

	render(): JSX.Element {
		return <div>
			<nav className="pt-navbar">
				<div className="pt-navbar-group pt-align-left">
					<div className="pt-navbar-heading">AutoABS</div>
					<Loading size="small"/>
				</div>
				<div className="pt-navbar-group pt-align-right">
					<button
						className="pt-button pt-minimal pt-icon-refresh"
						onClick={this.sync}
					/>
				</div>
			</nav>
			<InfiniteFlex
				style={css.builds}
				width={260}
				height={123}
				padding={5}
				scrollMargin={2}
				scrollMarginHit={1}
				buildItem={this.buildItem}
				traverse={this.traverse}
				items={this.state.builds}
				index={this.state.index}
				count={this.state.count}
			/>
			<BuildInfo/>
		</div>;
	}
}
