#!/usr/bin/env python3

import rclpy
import json
import time
from rclpy.node import Node

def get_ros_graph_info(node: Node) -> dict:
    graph_info = {
        "nodes": [],
        "topics": {}
    }

    time.sleep(1)

    # get all node names and namespaces
    node_names_and_namespaces = node.get_node_names_and_namespaces()
    graph_info["nodes"] = [
        f"{namespace}{name}" for name, namespace in node_names_and_namespaces
    ]

    # get all topics and their connections
    topic_names_and_types = node.get_topic_names_and_types(no_demangle=False)

    for topic_name, message_types in topic_names_and_types:
        # publishers
        publishers_info = node.get_publishers_info_by_topic(topic_name)
        publisher_nodes = [
            f"{info.node_namespace}{info.node_name}" for info in publishers_info
        ]

        # subscribers
        subscriptions_info = node.get_subscriptions_info_by_topic(topic_name)
        subscriber_nodes = [
            f"{info.node_namespace}{info.node_name}" for info in subscriptions_info
        ]

        if publisher_nodes or subscriber_nodes:
            graph_info["topics"][topic_name] = {
                "types": message_types,
                "publishers": publisher_nodes,
                "subscribers": subscriber_nodes,
            }

    return graph_info

def main(args=None):
    print("--- ROS 2 Graph Inspector ---")
    print("Initializing ROS 2 node to inspect the graph...")

    try:
        rclpy.init(args=args)
        inspector_node = Node("ros_graph_inspector_node")

        graph_data = get_ros_graph_info(inspector_node)
        json_output = json.dumps(graph_data, indent=2)

        print("\n--- Discovered ROS 2 Graph (JSON) ---")
        print(json_output)
        print("-------------------------------------\n")

    except Exception as e:
        print(f"An error occurred: {e}")
    finally:
        if rclpy.ok():
            print("Shutting down ROS 2 node.")
            inspector_node.destroy_node()
            rclpy.shutdown()

if __name__ == '__main__':
    main()
