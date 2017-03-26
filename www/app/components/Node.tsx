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
		height: '130px',
		margin: '5px',
		padding: '0',
	} as React.CSSProperties,
	content: {
		padding: '10px 0 0 10px',
	} as React.CSSProperties,
	name: {
		fontSize: '14px',
	} as React.CSSProperties,
	info: {
		fontSize: '12px',
		marginTop: '5px',
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

		return <div className="pt-card" style={css.card}>
			<div className="layout horizontal">
				<div style={css.content} className="card-content flex">
					<div className="layout vertical">
						<div style={css.name}>{node.id}</div>
						<div className="pt-text-muted" style={css.info}>
							{node.memory}
						</div>
					</div>
				</div>
			</div>
		</div>;
	}
}
