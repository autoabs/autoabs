/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as NodeTypes from '../types/NodeTypes';
import * as Blueprint from '@blueprintjs/core';

type OnClose = () => void;

interface Props {
	node: NodeTypes.Node;
	open: boolean;
	onClose: OnClose;
}

const css = {
	dialog: {
		maxWidth: 'calc(100% - 40px)',
	} as React.CSSProperties,
};

export default class NodeSettings extends React.Component<Props, void> {
	builderSettings(): JSX.Element {
		return <div>
			<div className="pt-text-muted">
				concurrency: {this.props.node.settings.concurrency}
			</div>
		</div>;
	}

	render(): JSX.Element {
		let node = this.props.node;
		let settings: JSX.Element;

		switch (node.type) {
			case 'builder':
				settings = this.builderSettings();
				break;
		}

		return <Blueprint.Dialog
				title={node.id}
				style={css.dialog}
				isOpen={this.props.open}
				onClose={this.props.onClose}
				canOutsideClickClose={false}
			>
				<div className="pt-dialog-body">
					<div className="pt-text-muted">
						id: {node.id}
					</div>
					<div className="pt-text-muted">
						type: {node.type}
					</div>
					{settings}
				</div>
				<div className="pt-dialog-footer">
					<div className="pt-dialog-footer-actions">
						<button type="button"
							className="pt-button"
							onClick={this.props.onClose}
						>Close</button>
					</div>
				</div>
		</Blueprint.Dialog>;
	}
}
