/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as NodeTypes from '../types/NodeTypes';
import * as Blueprint from '@blueprintjs/core';
import NodeSettings from './NodeSettings';

interface Props {
	node: NodeTypes.Node;
}

interface State {
	settings: boolean;
}

const css = {
	card: {
		flex: '1 0 auto',
		minWidth: '270px',
		maxWidth: '600px',
		margin: '5px',
		padding: '0',
	} as React.CSSProperties,
	layout: {
		position: 'relative',
	} as React.CSSProperties,
	content: {
		padding: '10px',
	} as React.CSSProperties,
	name: {
		fontSize: '14px',
	} as React.CSSProperties,
	info: {
		fontSize: '12px',
		marginTop: '5px',
	} as React.CSSProperties,
	settings: {
		position: 'absolute',
		right: '0',
		margin: '10px 10px 0 0',
	} as React.CSSProperties,
	stat: {
		marginTop: '1px',
	} as React.CSSProperties,
	settingsDialog: {
		maxWidth: 'calc(100% - 40px)',
	} as React.CSSProperties,
};

export default class Node extends React.Component<Props, State> {
	constructor(props: Props, context: any) {
		super(props, context);
		this.state = {
			settings: false,
		};
	}

	openDialog = (): void => {
		this.setState({
			settings: true,
		});
	}

	closeDialog = (): void => {
		this.setState({
			settings: false,
		});
	}

	render(): JSX.Element {
		let node = this.props.node;

		let memoryStyle = {
			width: node.memory + '%',
		} as React.CSSProperties;
		let load1Style = {
			width: node.load1 + '%',
		} as React.CSSProperties;
		let load5Style = {
			width: node.load5 + '%',
		} as React.CSSProperties;
		let load15Style = {
			width: node.load15 + '%',
		} as React.CSSProperties;

		return <div className="pt-card" style={css.card}>
			<div className="layout horizontal" style={css.layout}>
				<div style={css.content} className="card-content flex">
					<div className="layout vertical">
						<div style={css.name}>
							{node.id} - {node.type}
						</div>
						<div className="pt-text-muted" style={css.info}>
							memory:
						</div>
						<div style={css.stat}
							className="pt-progress-bar pt-no-stripes pt-intent-primary"
						>
							<div className="pt-progress-meter" style={memoryStyle}/>
						</div>
						<div className="pt-text-muted" style={css.info}>
							load1:
						</div>
						<div style={css.stat}
							className="pt-progress-bar pt-no-stripes pt-intent-success"
						>
							<div className="pt-progress-meter" style={load1Style}/>
						</div>
						<div className="pt-text-muted" style={css.info}>
							load5:
						</div>
						<div style={css.stat}
							className="pt-progress-bar pt-no-stripes pt-intent-warning"
						>
							<div className="pt-progress-meter" style={load5Style}/>
						</div>
						<div className="pt-text-muted" style={css.info}>
							load15:
						</div>
						<div style={css.stat}
							className="pt-progress-bar pt-no-stripes pt-intent-danger"
						>
							<div className="pt-progress-meter" style={load15Style}/>
						</div>
					</div>
				</div>
				<div>
					<button type="button"
						className="pt-button pt-minimal pt-icon-cog"
						style={css.settings}
						onClick={this.openDialog}
					/>
				</div>
			</div>
			<NodeSettings
				node={node}
				open={this.state.settings}
				onClose={this.closeDialog}
			/>
		</div>;
	}
}
