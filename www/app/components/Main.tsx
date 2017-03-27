/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as BuildActions from '../actions/BuildActions';
import * as NodeActions from '../actions/NodeActions';
import FixContainer from './FixContainer';
import Loading from './Loading';
import Builds from './Builds';
import Nodes from './Nodes';

document.body.className = 'root pt-dark';

interface State {
	page: string;
}

const css = {
	nav: {
		overflowX: 'auto',
		overflowY: 'hidden',
	} as React.CSSProperties,
	heading: {
		marginRight: '11px',
	} as React.CSSProperties,
}

export default class Main extends React.Component<void, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			page: 'builds',
		};
		this.sync();
	}

	setPage(page: string): void {
		this.setState({
			page: page,
		}, this.sync);
	}

	sync = (): void => {
		switch (this.state.page) {
			case 'builds':
				BuildActions.sync();
				break;
			case 'nodes':
				NodeActions.sync();
				break;
		}
	}

	render(): JSX.Element {
		let page: JSX.Element;

		switch (this.state.page) {
			case 'builds':
				page = <Builds/>;
				break;
			case 'nodes':
				page = <Nodes/>;
				break;
		}

		return <FixContainer>
			<nav className="pt-navbar layout horizontal" style={css.nav}>
				<div className="pt-navbar-group pt-align-left flex">
					<div className="pt-navbar-heading"
						style={css.heading}
					>AutoABS</div>
					<Loading size="small"/>
				</div>
				<div className="pt-navbar-group pt-align-right">
					<button
						className="pt-button pt-minimal pt-icon-oil-field"
						onClick={() => {this.setPage('builds')}}
					>Builds</button>
					<button
						className="pt-button pt-minimal pt-icon-layers"
						onClick={() => {this.setPage('nodes')}}
					>Nodes</button>
					<button
						className="pt-button pt-minimal pt-icon-refresh"
						onClick={this.sync}
					>Refresh</button>
				</div>
			</nav>
			{page}
		</FixContainer>;
	}
}
