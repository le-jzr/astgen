
#include <memory>
#include <vector>

class Visitor;

class NodeBase {

public:
	virtual NodeBase *clone() const = 0;
	virtual void visit(Visitor &visitor) = 0;
};
