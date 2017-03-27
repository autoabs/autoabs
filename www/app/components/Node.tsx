/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as NodeTypes from '../types/NodeTypes';

interface Props {
	node: NodeTypes.Node;
}

interface State {
	locked: boolean;
}

const css = {
	card: {
		flex: '1 0 auto',
		minWidth: '270px',
		maxWidth: '600px',
		margin: '5px',
		padding: '0',
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
	stat: {
		marginTop: '1px',
	} as React.CSSProperties,
};

export default class Node extends React.Component<Props, State> {
	constructor(props: Props, context: any) {
		super(props, context);
		this.state = {
			locked: false,
		};
	}

	render(): JSX.Element {
		let node = this.props.node;

		let memoryStyle = {
			'width': node.memory + '%',
		} as React.CSSProperties;
		let load1Style = {
			'width': node.load1 + '%',
		} as React.CSSProperties;
		let load5Style = {
			'width': node.load5 + '%',
		} as React.CSSProperties;
		let load15Style = {
			'width': node.load15 + '%',
		} as React.CSSProperties;

		return <div className="pt-card" style={css.card}>
			<div className="layout horizontal">
				<div style={css.content} className="card-content flex">
					<div className="layout vertical">
						<div style={css.name}>
							{node.id} - {node.type}
						</div>
						<div className="pt-text-muted" style={css.info}>
							memory:
						</div>
						<div style={css.stat}
							className="pt-progress-bar pt-no-stripes"
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
							className="pt-progress-bar pt-no-stripes pt-intent-primary"
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
			</div>
		</div>;
	}
}
