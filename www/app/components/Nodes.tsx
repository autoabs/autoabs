/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as NodeActions from '../actions/NodeActions';
import * as NodeTypes from '../types/NodeTypes';
import NodeStore from '../stores/NodeStore';
import Node from './Node';

interface State {
	nodes: NodeTypes.Nodes;
}

function getState(): State {
	return {
		nodes: NodeStore.nodes,
	};
}

const css = {
	nodes: {
		width: '100%',
		flex: 1,
		display: 'flex',
		flexDirection: 'row',
		flexWrap: 'wrap',
		padding: '5px',
	} as React.CSSProperties,
};

export default class Nodes extends React.Component<null, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = getState();
		NodeActions.sync();
	}

	componentDidMount(): void {
		NodeStore.addChangeListener(this.onChange);
	}

	componentWillUnmount(): void {
		NodeStore.removeChangeListener(this.onChange);
	}

	onChange = (): void => {
		this.setState(getState());
	}

	sync = (): void => {
		NodeActions.sync();
	}

	render(): JSX.Element {
		let nodes: JSX.Element[] = [];

		if (this.state.nodes) {
			for (let node of this.state.nodes) {
				nodes.push(<Node key={node.id} node={node}/>);
			}
		}

		return <div style={css.nodes}>
			{nodes}
		</div>;
	}
}
