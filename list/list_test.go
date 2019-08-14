package list_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/susruth/immutable/list"
)

var _ = Describe("Immutable List", func() {
	Context("when appending elements to a list", func() {
		It("should create new versions of the list when using ints", func() {
			list := New(1, 2, 3, 4)
			Expect(list.Values(0)).Should(Equal([]interface{}{1, 2, 3, 4}))

			list.Append(5)
			Expect(list.Values(1)).Should(Equal([]interface{}{1, 2, 3, 4, 5}))
		})

		It("should create new versions of the list when using strings", func() {
			list := New("hi", "hello", "bye")
			Expect(list.Values(0)).Should(Equal([]interface{}{"hi", "hello", "bye"}))

			list.Append("good")
			Expect(list.Values(1)).Should(Equal([]interface{}{"hi", "hello", "bye", "good"}))
		})
	})

	Context("when updating elements on a list", func() {
		It("should create new versions of the list when using ints", func() {
			list := New(1, 2, 3, 4)
			Expect(list.Values(0)).Should(Equal([]interface{}{1, 2, 3, 4}))

			list.Update(2, 5)
			Expect(list.Values(1)).Should(Equal([]interface{}{1, 2, 5, 4}))
		})

		It("should create new versions of the list when using strings", func() {
			list := New("hi", "hello", "bye")
			Expect(list.Values(0)).Should(Equal([]interface{}{"hi", "hello", "bye"}))

			list.Update(2, "good")
			Expect(list.Values(1)).Should(Equal([]interface{}{"hi", "hello", "good"}))
		})
	})

	Context("when removing elements from a list", func() {
		It("should create new versions of the list when using ints", func() {
			list := New(1, 2, 3, 4)
			Expect(list.Values(0)).Should(Equal([]interface{}{1, 2, 3, 4}))

			list.Delete(2)
			Expect(list.Values(1)).Should(Equal([]interface{}{1, 2, 4}))
		})

		It("should create new versions of the list when using strings", func() {
			list := New("hi", "hello", "bye")
			Expect(list.Values(0)).Should(Equal([]interface{}{"hi", "hello", "bye"}))

			list.Delete(1)
			Expect(list.Values(1)).Should(Equal([]interface{}{"hi", "bye"}))
		})
	})
})
