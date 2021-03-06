/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as NodeTypes from '../types/NodeTypes';
import * as NodeActions from '../actions/NodeActions';
import * as Blueprint from '@blueprintjs/core';

type OnClose = () => void;

interface Props {
	node: NodeTypes.Node;
	open: boolean;
	onClose: OnClose;
}

interface State {
	settings: NodeTypes.NodeSettings;
}

const css = {
	dialog: {
		maxWidth: 'calc(100% - 40px)',
	} as React.CSSProperties,
};

export default class NodeSettings extends React.Component<Props, State> {
	constructor(props: Props, context: any) {
		super(props, context);
		this.state = {
			settings: null,
		};
	}

	concurrencyChange = (val: number): void => {
		this.setState({
			settings: {
				...this.state.settings,
				concurrency: val,
			},
		});
	}

	builderSettings(): JSX.Element {
		let concurrency: number;

		if (this.state.settings) {
			let settings = this.state.settings as NodeTypes.BuilderSettings;
			concurrency = settings.concurrency;
		} else {
			let settings = this.props.node.settings as NodeTypes.BuilderSettings;
			concurrency = settings.concurrency;
		}

		return <div>
			<label className="pt-label">
				Concurrency
				<Blueprint.Slider
					min={1}
					max={10}
					stepSize={1}
					value={concurrency}
					onChange={this.concurrencyChange}
				/>
			</label>
		</div>;
	}

	onSave = (): void => {
		if (!this.state.settings) {
			this.onClose();
			return;
		}

		NodeActions.commit(this.props.node.id, this.state.settings).then(() => {
			this.props.onClose();
			setTimeout(() => {
				this.setState({
					settings: null,
				});
			}, 100);
		});
	}

	onClose = (): void => {
		this.props.onClose();
		setTimeout(() => {
			this.setState({
				settings: null,
			});
		}, 100);
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
				onClose={this.onClose}
				canOutsideClickClose={false}
			>
				<div className="pt-dialog-body">
					<label className="pt-label">
						Node ID
						<div className="pt-text-muted">
							{node.id}
						</div>
					</label>
					<label className="pt-label">
						Node Type
						<div className="pt-text-muted">
							{node.type}
						</div>
					</label>
					{settings}
				</div>
				<div className="pt-dialog-footer">
					<div className="pt-dialog-footer-actions">
						<button type="button"
							className="pt-button"
							onClick={this.onClose}
						>Close</button>
						<button type="button"
							className="pt-button pt-intent-primary"
							onClick={this.onSave}
						>Save</button>
					</div>
				</div>
		</Blueprint.Dialog>;
	}
}
