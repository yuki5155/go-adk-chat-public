import json
import ast
import operator

_ALLOWED_OPS = {
    ast.Add: operator.add,
    ast.Sub: operator.sub,
    ast.Mult: operator.mul,
    ast.Div: operator.truediv,
    ast.Pow: operator.pow,
    ast.USub: operator.neg,
    ast.UAdd: operator.pos,
    ast.Mod: operator.mod,
}


def _eval(node):
    if isinstance(node, ast.Constant):
        if not isinstance(node.value, (int, float)):
            raise ValueError("Only numeric literals are allowed")
        return node.value
    if isinstance(node, ast.BinOp):
        op = type(node.op)
        if op not in _ALLOWED_OPS:
            raise ValueError(f"Unsupported operator: {op.__name__}")
        return _ALLOWED_OPS[op](_eval(node.left), _eval(node.right))
    if isinstance(node, ast.UnaryOp):
        op = type(node.op)
        if op not in _ALLOWED_OPS:
            raise ValueError(f"Unsupported operator: {op.__name__}")
        return _ALLOWED_OPS[op](_eval(node.operand))
    raise ValueError(f"Unsupported expression type: {type(node).__name__}")


def _parse_body(event: dict) -> dict:
    """Support both API Gateway proxy format and direct RIE invocation."""
    if "body" in event:
        raw = event["body"]
        return json.loads(raw) if isinstance(raw, str) else raw
    return event


def lambda_handler(event, context):
    try:
        body = _parse_body(event)
        expression = body.get("expression", "")
        if not expression:
            return _response(400, {"error": "expression is required"})

        tree = ast.parse(expression, mode="eval")
        result = _eval(tree.body)

        return _response(200, {"result": result, "expression": expression})
    except (ValueError, SyntaxError) as e:
        return _response(400, {"error": str(e)})
    except Exception as e:
        return _response(500, {"error": "internal error"})


def _response(status: int, body: dict) -> dict:
    return {
        "statusCode": status,
        "headers": {"Content-Type": "application/json"},
        "body": json.dumps(body),
    }
